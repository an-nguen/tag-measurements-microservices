package ru.tagmeasurements.fetch_service.models;

import lombok.Data;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;


@Data
public class CloudHttpClient {
  private static final Logger log = LoggerFactory.getLogger(CloudHttpClient.class);
  private String sessionId;
  private WirelessTagAccount account;
  private List<TagManager> tagManagerList;
  private List<Tag> tagList;
  private List<Measurement> measurements;

  public CloudHttpClient(WirelessTagAccount account) {
    this.account = account;
  }


  public List<TagManager> getTagManagerList() {
    return tagManagerList;
  }

  public void setTagManagerList(List<TagManager> tagManagerList) {
    this.tagManagerList = tagManagerList;
  }
}
