package ru.tagmeasurements.fetch_service.services;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.google.gson.Gson;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpHeaders;
import org.springframework.stereotype.Service;
import ru.tagmeasurements.fetch_service.dtos.*;
import ru.tagmeasurements.fetch_service.models.Measurement;
import ru.tagmeasurements.fetch_service.models.Tag;
import ru.tagmeasurements.fetch_service.models.TagManager;
import ru.tagmeasurements.fetch_service.models.WirelessTagAccount;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.time.temporal.ChronoUnit;
import java.util.ArrayList;
import java.util.Collections;
import java.util.Comparator;
import java.util.List;
import java.util.stream.Collectors;

@Service
public class ApiWrapperService {
    private final HttpClientService client;
    private final Gson gson;
    private final Logger log = LoggerFactory.getLogger(ApiWrapperService.class);

    @Autowired
    public ApiWrapperService(HttpClientService client, Gson gson) {
        this.client = client;
        this.gson = gson;

    }

    private HttpHeaders getHttpHeaders(String sessionId) {
        var headers = new HttpHeaders();
        headers.add("Content-Type", "application/json");
        headers.add("Accept", "*/*");
        headers.add("Cookie", sessionId);
        return headers;
    }

    private String selectTagManager(String mac, String sessionId) {
        var response = client
                .post("/ethAccount.asmx/SelectTagManager",
                        gson.toJson(new SelectTagManagerRequest(mac)), getHttpHeaders(sessionId));
        return String.join("", response.getHeaders().get("Set-Cookie"));
    }

    public String getSessionId(WirelessTagAccount account) {
        var response = client.post("/ethAccount.asmx/SignIn", gson.toJson(account));
        return String.join("", response.getHeaders().get("Set-Cookie"));
    }

    private TempSensorConfig getTagConfig(String sessionId, Tag tag) {
        var resp = client.post("/ethClient.asmx/LoadTempSensorConfig",
                "{id: " + tag.getSlaveId() + "}",
                getHttpHeaders(sessionId));
        return gson.fromJson(resp.getBody(), LoadTempSensorConfigResponse.class).getD();
    }

    public List<TagManager> getTagManagers(String sessionId, WirelessTagAccount account) {
        var result = new ArrayList<TagManager>();
        var response = client.post("/ethAccount.asmx/GetTagManagers",
                "{}", getHttpHeaders(sessionId));
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
            tagManager.setEmail(account.getEmail());
            result.add(tagManager);
        }
        return result;
    }

    public List<Tag> getTags(String sessionId, List<TagManager> cloudTagManagers) {
        var result = new ArrayList<Tag>();
        for (var tm: cloudTagManagers) {
            var sessionIdSelected = selectTagManager(tm.getMac().replaceAll(":", "").toUpperCase(), sessionId);
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

    private GetTagListResponse getTagListResponse(String sessionId) {
        var response = client.post("/ethClient.asmx/GetTagList",
                "{}", getHttpHeaders(sessionId));
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

    public List<Measurement> getMeasurements(List<TagManager> tagManagers, String type) {
        var result = new ArrayList<Measurement>();
        for (var tm: tagManagers) {
            var cloudTags = getTags(tm);
            var tags = cloudTags.stream()
                    .map(Tag::getSlaveId)
                    .sorted(Comparator.naturalOrder())
                    .collect(Collectors.toList());
            var request = new GetMultiTagStatsRawRequest();
            var responses = new ArrayList<GetMultiTagStatsRawResponse>();
            for (var tag : tags) {
                try {
                    request.setIds(Collections.singletonList(tag));
                    request.setFromDate(LocalDate.now().minus(75, ChronoUnit.DAYS));
                    request.setToDate(LocalDate.now());
                    request.setType(type);
                    var jsonRequest = gson.toJson(request);
                    var response = client.post("/ethLogs.asmx/GetMultiTagStatsRaw",
                            jsonRequest, getHttpHeaders(tm.getSessionId()));
                    if (response.getStatusCode().is2xxSuccessful()) {
                        responses.add(gson.fromJson(response.getBody(), GetMultiTagStatsRawResponse.class));
                    }
                } catch (Exception e) {
                    log.error(e.getLocalizedMessage());
                }
            }
            for (var json : responses) {
                var jsonResponse = json.getD();
                for (var stat : jsonResponse.getStats()) {
                    var date = LocalDate.parse(stat.getDate(), DateTimeFormatter.ofPattern("M/d/yyyy"));
                    for (var j = 0; j < stat.values.size(); j++) {
                        var measurement = new Measurement();
                        measurement
                                .setDate(LocalDateTime
                                        .parse(date.format(DateTimeFormatter.ISO_DATE) + "T00:00:00")
                                        .plus(stat.tods.get(0).get(j), ChronoUnit.SECONDS));
                        switch (type) {
                            case "temperature" -> measurement.setTemperature(stat.values.get(0).get(j));
                            case "cap" -> measurement.setHumidity(stat.values.get(0).get(j));
                            case "batteryVolt" -> measurement.setVoltage(stat.values.get(0).get(j));
                            case "signal" -> measurement.setSignal(stat.values.get(0).get(j));
                        }
                        var foundTag = cloudTags.stream().filter(tag -> tag.getSlaveId().equals(stat.getIds().get(0))).findFirst();
                        foundTag.ifPresent(tag -> measurement.setTagUUID(tag.getUuid()));
                        result.add(measurement);
                    }

                }
            }
        }
        return result;
    }
}
