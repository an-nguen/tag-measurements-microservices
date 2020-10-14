package ru.tagmeasurements.fetch_service.utils;

import org.springframework.http.HttpHeaders;

public class HttpHelpers {

  public static HttpHeaders getHttpHeaders(String sessionId) {
    var headers = new HttpHeaders();
    headers.add("Content-Type", "application/json");
    headers.add("Accept", "*/*");
    headers.add("Cookie", sessionId);
    return headers;
  }
}
