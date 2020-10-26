package ru.tagmeasurements.fetch_service.services;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Service;
import ru.tagmeasurements.fetch_service.models.CloudHttpClient;

@Service
public class ScheduledService {
  private final Logger log = LoggerFactory.getLogger(ScheduledService.class);

  private final TagService tagService;
  private final TagManagerService tagManagerService;
  private final MeasurementService measurementService;
  private final MeasurementRTService measurementRTService;
  private final WirelessTagAccountService wirelessTagAccountService;
  private final String[] types = {
    "temperature", "cap", "batteryVolt", "signal"
  };

  @Autowired
  public ScheduledService(TagService tagService,
                          TagManagerService tagManagerService,
                          MeasurementService measurementService,
                          MeasurementRTService measurementRTService, WirelessTagAccountService wirelessTagAccountService) {
    this.tagService = tagService;
    this.tagManagerService = tagManagerService;
    this.measurementService = measurementService;
    this.measurementRTService = measurementRTService;
    this.wirelessTagAccountService = wirelessTagAccountService;
  }
  @Scheduled(fixedDelay = 60000)
  private void clientFetchMeasurements() {
    var accounts = wirelessTagAccountService.findAll();
    for (var account: accounts) {
      var client = new CloudHttpClient(account);
      client.setSessionId(tagManagerService.getSessionId(client.getAccount()));
      client.setTagManagerList(tagManagerService.getTagManagers(client));
      client.setTagList(tagService.getTags(client));
      measurementService.storeMeasurements(client.getTagManagerList(), types);
    }
  }

  @Scheduled(fixedDelay = 300000)
  private void clientFetchTagsAndTagManagers() {
    var accounts = wirelessTagAccountService.findAll();
    for (var account: accounts) {
      var client = new CloudHttpClient(account);
      client.setSessionId(tagManagerService.getSessionId(client.getAccount()));
      client.setTagManagerList(tagManagerService.getTagManagers(client));
      client.setTagList(tagService.getTags(client));
      tagService.storeTags(client, tagManagerService.findAllByEmail(client.getAccount().getEmail()));
      tagManagerService.storeTagManagers(client);
    }
  }

//  @Scheduled(fixedDelay = 40000)
//  private void clientFetchMeasurementsRT() {
//    var accounts = wirelessTagAccountService.findAll();
//    for (var account: accounts) {
//      var client = new CloudHttpClient(account);
//      client.setSessionId(tagManagerService.getSessionId(client.getAccount()));
//      client.setTagManagerList(tagManagerService.getTagManagers(client));
//      measurementRTService.storeMeasurements(client);
//    }
//  }
}
