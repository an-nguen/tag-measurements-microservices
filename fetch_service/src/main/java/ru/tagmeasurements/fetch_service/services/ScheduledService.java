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
  private final WirelessTagAccountService wirelessTagAccountService;
  private final String[] types = {
    "temperature", "cap", "batteryVolt", "signal"
  };

  @Autowired
  public ScheduledService(TagService tagService,
                          TagManagerService tagManagerService,
                          MeasurementService measurementService,
                          WirelessTagAccountService wirelessTagAccountService) {
    this.tagService = tagService;
    this.tagManagerService = tagManagerService;
    this.measurementService = measurementService;
    this.wirelessTagAccountService = wirelessTagAccountService;
  }
  @Scheduled(fixedDelay = 15000)
  private void clientFetchMeasurements() {
    log.info("Begin client fetch measurements");
    var accounts = wirelessTagAccountService.findAll();
    for (var account: accounts) {
      var client = new CloudHttpClient(account);
      client.setSessionId(tagManagerService.getSessionId(client.getAccount()));
      client.setTagManagerList(tagManagerService.getTagManagers(client));
      client.setTagList(tagService.getTags(client));
      measurementService.storeMeasurements(client.getTagManagerList(), types);
    }
    log.info("End client fetch measurements");
  }

  @Scheduled(fixedDelay = 15000)
  private void clientFetchTagsAndTagManagers() {
    log.info("Begin client fetch tags and tag managers");
    var accounts = wirelessTagAccountService.findAll();
    for (var account: accounts) {
      var client = new CloudHttpClient(account);
      client.setSessionId(tagManagerService.getSessionId(client.getAccount()));
      client.setTagManagerList(tagManagerService.getTagManagers(client));
      client.setTagList(tagService.getTags(client));
      tagService.storeTags(client, tagManagerService.findAllByEmail(client.getAccount().getEmail()));
      tagManagerService.storeTagManagers(client);
    }
    log.info("End client fetch tags and tag managers");
  }

}
