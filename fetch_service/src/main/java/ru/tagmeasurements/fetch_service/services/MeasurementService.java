package ru.tagmeasurements.fetch_service.services;

import com.google.gson.Gson;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import ru.tagmeasurements.fetch_service.dtos.GetMultiTagStatsRawRequest;
import ru.tagmeasurements.fetch_service.dtos.GetMultiTagStatsRawResponse;
import ru.tagmeasurements.fetch_service.dtos.MultiTagStats;
import ru.tagmeasurements.fetch_service.dtos.MultiTagStatsRawResponse;
import ru.tagmeasurements.fetch_service.models.Measurement;
import ru.tagmeasurements.fetch_service.models.Tag;
import ru.tagmeasurements.fetch_service.models.TagManager;
import ru.tagmeasurements.fetch_service.repositories.MeasurementRepository;
import ru.tagmeasurements.fetch_service.utils.HttpHelpers;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.time.temporal.ChronoUnit;
import java.util.*;
import java.util.stream.Collectors;

@Service
public class MeasurementService {
  private final Logger log = LoggerFactory.getLogger(MeasurementService.class);

  private final MeasurementRepository repository;
  private final HttpClientService client;
  private final TagService tagService;
  private final Gson gson;

  @Autowired
  public MeasurementService(MeasurementRepository repository, HttpClientService client, TagService tagService, Gson gson) {
    this.repository = repository;
    this.client = client;
    this.tagService = tagService;
    this.gson = gson;
  }


  public void storeMeasurements(List<TagManager> tagManagers, String[] types) {
    var result = new HashMap<LocalDateTime, List<Measurement>>();
    for (var type: types) {
      for (var tm : tagManagers) {
        var cloudTags = tagService.getTags(tm).stream()
          .sorted(Comparator.comparingLong(Tag::getSlaveId))
          .collect(Collectors.toList());
        var tagSlaveIds = cloudTags.stream()
          .map(Tag::getSlaveId)
          .collect(Collectors.toList());
        var request = new GetMultiTagStatsRawRequest();
        // Fetch measurement of one tag
        for (var tagSlaveId: tagSlaveIds) {
          try {
            var optionalCloudTag = cloudTags.stream().filter(ct -> ct.getSlaveId().equals(tagSlaveId)).findFirst();
            if(optionalCloudTag.isPresent()) {
              var cloudTag = optionalCloudTag.get();
              // Send request
              request.setIds(Collections.singletonList(tagSlaveId));
              // Latest added record by uuid
//              var latestRecord = repository.getFirstByTagUUIDOrderByDateDesc(cloudTag.getUuid());
              var fromDate = LocalDateTime.now().minus(75, ChronoUnit.DAYS);
              request.setFromDate(fromDate.toLocalDate());
              request.setToDate(LocalDate.now());
              request.setType(type);
              var response = getMeasurements(tm, request);
              // Handle request
              handleResponse(response, type, cloudTag, result);
            } else {
              log.info("can't find cloud tag by slave id");
            }
          } catch (Exception e) {
            log.error(e.getLocalizedMessage());
          }
        }
      }
    }
    // Store to database
    if (result.size() > 0) {
      int count = 0;
      for (var set: result.entrySet()) {
        repository.saveAll(set.getValue());
        count++;
      }
      log.info(String.format("The %d measurements stored.", count));
    }
  }

  private void handleResponse(MultiTagStatsRawResponse response, String type, Tag cloudTag, HashMap<LocalDateTime, List<Measurement>> result) {
    for (var stat : response.getStats()) {
      // Offset date
      var date = LocalDate.parse(stat.getDate(), DateTimeFormatter.ofPattern("M/d/yyyy"));
      // Iterate through values/tods
      for (var j = 0; j < stat.values.get(0).size(); j++) {
        // Set date time record
        var dateTime = LocalDateTime
          .parse(date.format(DateTimeFormatter.ISO_DATE) + "T00:00:00")
          .plus(stat.tods.get(0).get(j), ChronoUnit.SECONDS);
        // Attempt to find similar record by date and uuid in result collection
        if (result.containsKey(dateTime)) {
          var list = result.get(dateTime);
          var optionalMeasurement = list.stream()
            .filter(measurement -> measurement.getTagUUID().equals(cloudTag.getUuid()))
            .findFirst();
          int finalJ = j;
          optionalMeasurement.ifPresentOrElse(measurement -> {
            // Update record
            int n = list.indexOf(measurement);
            if (n == -1) {
              return;
            }
            fillMeasurement(type, stat, finalJ, measurement);
            var databaseMeasurement = repository
              .getFirstByDateAndTagUUIDOrderByDateDesc(dateTime, cloudTag.getUuid());
            if (databaseMeasurement != null) {
              measurement.setId(databaseMeasurement.getId());
            }
            list.set(n, measurement);
          }, () -> {
            // New one record
            addNewMeasurement(type, cloudTag, stat, finalJ, dateTime, list);
          });
        } else {
          result.put(dateTime, new LinkedList<>());
          var list = result.get(dateTime);
          addNewMeasurement(type, cloudTag, stat, j, dateTime, list);
        }
      }
    }
  }

  private void addNewMeasurement(String type, Tag cloudTag, MultiTagStats stat, int j, LocalDateTime dateTime, List<Measurement> list) {
    var measurement = new Measurement();
    measurement
      .setDate(dateTime);
    fillMeasurement(type, stat, j, measurement);
    var databaseMeasurement = repository
      .getFirstByDateAndTagUUIDOrderByDateDesc(dateTime, cloudTag.getUuid());
    if (databaseMeasurement != null)
      measurement.setId(databaseMeasurement.getId());
    measurement.setTagUUID(cloudTag.getUuid());
    list.add(measurement);
  }

  private MultiTagStatsRawResponse getMeasurements(TagManager tm, GetMultiTagStatsRawRequest request) throws Exception {
    var jsonRequest = gson.toJson(request);
    var response = client.post("/ethLogs.asmx/GetMultiTagStatsRaw",
      jsonRequest, HttpHelpers.getHttpHeaders(tm.getSessionId()));
    if (!response.getStatusCode().is2xxSuccessful()) {
      throw new Exception(String.valueOf(response.getStatusCode()));
    }
    return gson.fromJson(response.getBody(), GetMultiTagStatsRawResponse.class).getD();

  }

  private void fillMeasurement(String type, MultiTagStats stat, int j, Measurement measurement) {
    switch (type) {
      case "temperature" -> measurement.setTemperature(stat.values.get(0).get(j));
      case "cap" -> measurement.setHumidity(stat.values.get(0).get(j));
      case "batteryVolt" -> measurement.setVoltage(stat.values.get(0).get(j));
      case "signal" -> measurement.setSignal(stat.values.get(0).get(j));
    }
  }
}
