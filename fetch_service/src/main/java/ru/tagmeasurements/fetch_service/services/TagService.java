package ru.tagmeasurements.fetch_service.services;

import com.google.gson.Gson;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import ru.tagmeasurements.fetch_service.dtos.GetTagListResponse;
import ru.tagmeasurements.fetch_service.dtos.LoadTempSensorConfigResponse;
import ru.tagmeasurements.fetch_service.dtos.TempSensorConfig;
import ru.tagmeasurements.fetch_service.models.CloudHttpClient;
import ru.tagmeasurements.fetch_service.models.Tag;
import ru.tagmeasurements.fetch_service.models.TagManager;
import ru.tagmeasurements.fetch_service.repositories.TagRepository;
import ru.tagmeasurements.fetch_service.utils.HttpHelpers;

import java.time.LocalDate;
import java.util.*;
import java.util.stream.Collectors;

@Service
public class TagService {
  private final Logger log = LoggerFactory.getLogger(TagService.class);
  private final HttpClientService client;
  private final Gson gson;

  private final TagRepository repository;

  @Autowired
  public TagService(HttpClientService client, Gson gson, TagRepository repository) {
    this.client = client;
    this.gson = gson;
    this.repository = repository;
  }

  public List<Tag> getTags(CloudHttpClient cloudClient) {
    var result = new ArrayList<Tag>();
    for (var tm: cloudClient.getTagManagerList()) {
      var response = client
        .post("/ethAccount.asmx/SelectTagManager",
          String.format("{\"mac\":\"%s\"}", tm.getMac().replace(":", "")), HttpHelpers.getHttpHeaders(cloudClient.getSessionId()));
      var sessionIdSelected = String.join("", response.getHeaders().get("Set-Cookie"));
      tm.setSessionId(sessionIdSelected);

      var jsonResponse = getTagListResponse(sessionIdSelected);
      for (var cloudTag: jsonResponse.getD()) {
        var tag = new Tag();
        tag.parse(cloudTag);
        tag.setVerificationDate(LocalDate.now());
        tag.parse(getTagConfig(sessionIdSelected, tag));
        tag.setMac(tm.getMac());
        result.add(tag);
      }
    }
    return result;
  }

  GetTagListResponse getTagListResponse(String sessionId) {
    var response = client.post("/ethClient.asmx/GetTagList",
      "{}", HttpHelpers.getHttpHeaders(sessionId));
    return gson.fromJson(response.getBody(), GetTagListResponse.class);
  }

  public List<Tag> getTags(TagManager tagManager) {
    var result = new ArrayList<Tag>();
    var jsonResponse = getTagListResponse(tagManager.getSessionId());
    for (var cloudTag: jsonResponse.getD()) {
      var tag = new Tag();
      tag.parse(cloudTag);
      tag.setVerificationDate(LocalDate.now());
      tag.parse(getTagConfig(tagManager.getSessionId(), tag));
      tag.setMac(tagManager.getMac());
      result.add(tag);
    }

    return result;
  }


  private TempSensorConfig getTagConfig(String sessionId, Tag tag) {
    var resp = client.post("/ethClient.asmx/LoadTempSensorConfig",
      "{id: " + tag.getSlaveId() + "}",
      HttpHelpers.getHttpHeaders(sessionId));
    return gson.fromJson(resp.getBody(), LoadTempSensorConfigResponse.class).getD();
  }

  public void storeTags(CloudHttpClient cloudClient, List<TagManager> databaseTagManagers) {
    var tagList = new ArrayList<Tag>();
    for (var databaseTagManager: databaseTagManagers) {
      var tags = repository.findAll()
        .stream()
        .filter(tag -> tag.getMac().equals(databaseTagManager.getMac())).collect(Collectors.toList());
      if (!tags.isEmpty()) {
        tagList.addAll(tags);
      }
    }
    var databaseTagSet = new HashSet<>(tagList);
    var cloudTagSet = new HashSet<>(cloudClient.getTagList());
    // tags to remove
    var removed = databaseTagSet.stream().filter(t -> {
      for (var tag: cloudTagSet) {
        if (t.getUuid().equals(tag.getUuid()))
          return false;
      }
      return true;
    }).collect(Collectors.toCollection(HashSet::new));
    if (removed.size() > 0) {
      log.info(String.format("Delete removed %d database tags", removed.size()));
      removed.forEach(t -> {
        repository.deleteById(t.getUuid());
      });
    }

    var databaseTags = repository.findAll();
    for (var tag: cloudTagSet) {
      var found = databaseTags.stream()
        .filter(databaseTag -> databaseTag.getUuid().equals(tag.getUuid()))
        .findFirst();
      found.ifPresentOrElse(t -> {
        tag.setVerificationDate(t.getVerificationDate());
      }, () -> {
        tag.setVerificationDate(LocalDate.now());
      });
    }

    if (cloudTagSet.size() > 0) {
//      log.info(String.format("Store %d tags to database", cloudTagSet.size()));
      repository.saveAll(cloudTagSet);
    }
  }
}
