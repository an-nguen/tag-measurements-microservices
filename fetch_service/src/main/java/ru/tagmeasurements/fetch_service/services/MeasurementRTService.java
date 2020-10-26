package ru.tagmeasurements.fetch_service.services;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import ru.tagmeasurements.fetch_service.models.CloudHttpClient;
import ru.tagmeasurements.fetch_service.models.MeasurementRT;
import ru.tagmeasurements.fetch_service.models.Tag;
import ru.tagmeasurements.fetch_service.repositories.MeasurementRTRepository;
import ru.tagmeasurements.fetch_service.utils.HttpHelpers;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.util.LinkedList;
import java.util.UUID;

@Service
public class MeasurementRTService {
  private final MeasurementRTRepository repository;
  private final TagService tagService;
  private final HttpClientService client;

  @Autowired
  public MeasurementRTService(MeasurementRTRepository repository, TagService tagService, HttpClientService client) {
    this.repository = repository;
    this.tagService = tagService;
    this.client = client;
  }

  public void storeMeasurements(CloudHttpClient cloudClient) {
    var result = new LinkedList<MeasurementRT>();
    for (var tm: cloudClient.getTagManagerList()) {
      var response = client
        .post("/ethAccount.asmx/SelectTagManager",
          String.format("{\"mac\":\"%s\"}", tm.getMac().replace(":", "")), HttpHelpers.getHttpHeaders(cloudClient.getSessionId()));
      var sessionIdSelected = String.join("", response.getHeaders().get("Set-Cookie"));
      tm.setSessionId(sessionIdSelected);

      var jsonResponse = tagService.getTagListResponse(sessionIdSelected);
      for (var cloudTag: jsonResponse.getD()) {
        var measurementRT = new MeasurementRT();
        measurementRT.setTagUUID(UUID.fromString(cloudTag.getUuid()));
        measurementRT.setTemperature(cloudTag.getTemperature());
        measurementRT.setHumidity(cloudTag.getCap());
        measurementRT.setVoltage(cloudTag.getBatteryVolt());
        measurementRT.setSignal((double) cloudTag.getSignaldBm());
        measurementRT.setDate(LocalDateTime.now());

        result.add(measurementRT);
      }
    }
    if (result.size() > 0) {
      repository.saveAll(result);
    }
  }
}
