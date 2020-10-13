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
import java.util.Comparator;
import java.util.LinkedList;
import java.util.List;
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
    var result = new LinkedList<Measurement>();
    for (var type: types) {
      for (var tm : tagManagers) {
        var cloudTags = tagService.getTags(tm).stream()
          .sorted(Comparator.comparingLong(Tag::getSlaveId))
          .collect(Collectors.toList());
        var tags = cloudTags.stream()
          .map(Tag::getSlaveId)
          .collect(Collectors.toList());
        var request = new GetMultiTagStatsRawRequest();
        try {
          request.setIds(tags);
          request.setFromDate(LocalDate.now().minus(75, ChronoUnit.DAYS));
          request.setToDate(LocalDate.now());
          request.setType(type);
          var response = getMeasurements(tm, request);

          for (var stat : response.getStats()) {
            var date = LocalDate.parse(stat.getDate(), DateTimeFormatter.ofPattern("M/d/yyyy"));
            for (int i = 0; i < cloudTags.size(); i++) {
              int finalI = i;
              for (var j = 0; j < stat.values.get(finalI).size(); j++) {
                var dateTime = LocalDateTime
                  .parse(date.format(DateTimeFormatter.ISO_DATE) + "T00:00:00")
                  .plus(stat.tods.get(finalI).get(j), ChronoUnit.SECONDS);
                var foundMeasurement = result.stream()
                  .filter(m -> m.getDate().equals(dateTime) && m.getTagUUID()
                    .equals(cloudTags.get(finalI).getUuid()))
                  .findFirst();
                int finalJ = j;
                foundMeasurement.ifPresentOrElse(measurement -> {
                  int n = result.indexOf(measurement);
                  if (n == -1) {
                    return;
                  }
                  fillMeasurement(type, stat, finalI, finalJ, measurement);
                  var databaseMeasurement = repository
                    .getFirstByDateAndTagUUIDOrderByDateDesc(dateTime, cloudTags.get(finalI).getUuid());
                  if (databaseMeasurement != null) {
                    measurement.setId(databaseMeasurement.getId());
                  }
                  result.set(n, measurement);
                }, () -> {
                  var measurement = new Measurement();
                  measurement
                    .setDate(dateTime);
                  fillMeasurement(type, stat, finalI, finalJ, measurement);
                  measurement.setTagUUID(cloudTags.get(finalI).getUuid());
                  result.add(measurement);
                });
              }
            }
          }
        } catch (Exception e) {
          log.error(e.getLocalizedMessage());
        }
      }
    }
    if (result.size() > 0) {
      repository.saveAll(result);
      log.info(String.format("The %d measurements stored.", result.size()));
    }
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

  private void fillMeasurement(String type, MultiTagStats stat, int i, int j, Measurement measurement) {
    switch (type) {
      case "temperature" -> measurement.setTemperature(stat.values.get(i).get(j));
      case "cap" -> measurement.setHumidity(stat.values.get(i).get(j));
      case "batteryVolt" -> measurement.setVoltage(stat.values.get(i).get(j));
      case "signal" -> measurement.setSignal(stat.values.get(i).get(j));
    }
  }
}
