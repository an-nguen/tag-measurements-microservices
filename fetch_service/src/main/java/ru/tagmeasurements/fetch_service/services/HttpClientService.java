package ru.tagmeasurements.fetch_service.services;

import org.springframework.http.*;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

@Service
public class HttpClientService {

    private final String host = "http://wirelesstag.net";
    private RestTemplate rest;
    private HttpHeaders httpHeaders;
    private HttpStatus status;

    public HttpClientService() {
        this.rest = new RestTemplate();
        this.httpHeaders = new HttpHeaders();
        httpHeaders.add("Content-Type", "application/json");
        httpHeaders.add("Accept", "*/*");
    }

    public ResponseEntity<String> get(String uri) {
        var requestEntity = new HttpEntity<>("", httpHeaders);
        var responseEntity = rest.exchange(host + uri, HttpMethod.GET, requestEntity, String.class);
        this.setStatus(responseEntity.getStatusCode());
        return responseEntity;
    }

    public ResponseEntity<String> post(String uri, String json) {
        var requestEntity = new HttpEntity<>(json, httpHeaders);
        var responseEntity = rest.exchange(host + uri, HttpMethod.POST, requestEntity, String.class);
        this.setStatus(responseEntity.getStatusCode());
        return responseEntity;
    }

    public ResponseEntity<String> post(String uri, String json, HttpHeaders customHeaders) {
        var reqEntity = new HttpEntity<>(json, customHeaders);
        var respEntity = rest.exchange(host + uri, HttpMethod.POST, reqEntity, String.class);
        this.setStatus(respEntity.getStatusCode());
        return respEntity;
    }

    public void put(String uri, String json) {
        var requestEntity = new HttpEntity<>(json, httpHeaders);
        var responseEntity = rest.exchange(host + uri, HttpMethod.PUT, requestEntity, (Class<Object>) null);
        this.setStatus(responseEntity.getStatusCode());
    }

    public void delete(String uri) {
        var requestEntity = new HttpEntity<>("", httpHeaders);
        var responseEntity = rest.exchange(host + uri, HttpMethod.DELETE, requestEntity, (Class<Object>) null);
        this.setStatus(responseEntity.getStatusCode());
    }

    public HttpStatus getStatus() {
        return status;
    }

    public void setStatus(HttpStatus status) {
        this.status = status;
    }
}
