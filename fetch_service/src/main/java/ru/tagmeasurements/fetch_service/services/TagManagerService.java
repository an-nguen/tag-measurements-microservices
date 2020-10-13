package ru.tagmeasurements.fetch_service.services;

import com.google.gson.Gson;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import ru.tagmeasurements.fetch_service.dtos.GetTagManagersResponse;
import ru.tagmeasurements.fetch_service.dtos.SelectTagManagerRequest;
import ru.tagmeasurements.fetch_service.models.CloudHttpClient;
import ru.tagmeasurements.fetch_service.models.TagManager;
import ru.tagmeasurements.fetch_service.models.WirelessTagAccount;
import ru.tagmeasurements.fetch_service.repositories.TagManagerRepository;
import ru.tagmeasurements.fetch_service.utils.HttpHelpers;

import java.util.ArrayList;
import java.util.List;

@Service
public class TagManagerService {
  private final Logger log = LoggerFactory.getLogger(TagManagerService.class);
  private final HttpClientService client;
  private final Gson gson;
  private final TagManagerRepository repository;

  @Autowired
  public TagManagerService(HttpClientService client, Gson gson, TagManagerRepository tagManagerRepository) {
    this.client = client;
    this.gson = gson;
    this.repository = tagManagerRepository;
  }

  public String getSessionId(WirelessTagAccount account) {
    var response = client.post("/ethAccount.asmx/SignIn", gson.toJson(account));
    return String.join("", response.getHeaders().get("Set-Cookie"));
  }

  public List<TagManager> findAllByEmail(String email) {
    return repository.findAllByEmail(email);
  }

  public List<TagManager> getTagManagers(CloudHttpClient cloudClient) {
    var result = new ArrayList<TagManager>();
    var response = client.post("/ethAccount.asmx/GetTagManagers",
      "{}", HttpHelpers.getHttpHeaders(cloudClient.getSessionId()));
    var jsonResponse = gson.fromJson(response.getBody(), GetTagManagersResponse.class);
    for (var tm: jsonResponse.getD()) {
      var tagManager = new TagManager();
      tagManager.parse(tm);
      var mac = tagManager.getMac().toLowerCase();
      StringBuilder stringBuilder = new StringBuilder();
      for (var i = 0; i < mac.length(); i++) {
        if ((i + 1) != mac.length() && ((i + 1) % 2 == 0)) {
          stringBuilder.append(mac.charAt(i));
          stringBuilder.append(":");
        } else {
          stringBuilder.append(mac.charAt(i));
        }
      }
      tagManager.setMac(stringBuilder.toString());
      tagManager.setEmail(cloudClient.getAccount().getEmail());
      result.add(tagManager);
    }
    return result;
  }

  public void storeTagManagers(CloudHttpClient httpClient) {
    log.info("Store fetched data to database");
    var databaseTagManagers = repository.findAllByEmail(httpClient.getAccount().getEmail());
    httpClient.getTagManagerList().removeAll(databaseTagManagers);
    if (httpClient.getTagManagerList().size() > 0) {
      log.info("Store tag managers to database");
      repository.saveAll(httpClient.getTagManagerList());
    }
  }
}
