package ru.tagmeasurements.fetch_service.models;

import lombok.Data;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import ru.tagmeasurements.fetch_service.repositories.MeasurementRepository;
import ru.tagmeasurements.fetch_service.repositories.TagManagerRepository;
import ru.tagmeasurements.fetch_service.repositories.TagRepository;
import ru.tagmeasurements.fetch_service.services.ApiWrapperService;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;


@Data
public class CloudHttpClient {
    private static final Logger log = LoggerFactory.getLogger(CloudHttpClient.class);
    private String sessionId;
    private WirelessTagAccount account;
    private ApiWrapperService service;
    private List<TagManager> tagManagerList;
    private List<Tag> tagList;
    private List<Measurement> measurements;

    public CloudHttpClient(WirelessTagAccount account, ApiWrapperService apiWrapperService) {
        this.account = account;
        this.service = apiWrapperService;
    }
    public void init() {
        try {
            this.sessionId = service.getSessionId(account);
            this.tagManagerList = service.getTagManagers(sessionId, account);
            this.tagList = service.getTags(sessionId, this.tagManagerList);
            this.measurements = service.getMeasurements(sessionId, "temperature", this.tagList);
        } catch (Exception e) {
            log.error(e.getLocalizedMessage());
        }
    }
    public void store(TagManagerRepository tagManagerRepository, TagRepository tagRepository, MeasurementRepository measurementRepository) {
        var tagManagers = tagManagerRepository.findAllByEmail(account.getEmail());
        this.tagManagerList.removeAll(tagManagers);
        if (this.tagManagerList.size() > 0) {
            tagManagerRepository.saveAll(this.tagManagerList);
        }
        var tagList = new ArrayList<Tag>();
        for (var tm: tagManagers) {
            var tempTags = tagRepository.findAll()
                    .stream()
                    .filter(tag -> tag.getMac().equals(tm.getMac())).collect(Collectors.toList());
            if (!tempTags.isEmpty()) {
                tagList.addAll(tempTags);
            }
        }
        this.tagList.removeAll(tagList);
        if (this.tagList.size() > 0) {
            tagRepository.saveAll(this.tagList);
        }
        if (this.measurements == null)  return;
        var measurementList = measurementRepository
                .getAllByDateInAndTagUUIDInOrderByDateDesc(
                        this.measurements.stream().map(Measurement::getDate).collect(Collectors.toSet()),
                        this.measurements.stream().map(Measurement::getTagUUID).collect(Collectors.toSet()));
        for (var measurement: this.measurements) {
            var databaseMeasurement = measurementList.stream()
                    .filter(m -> m.getDate().equals(measurement.getDate())
                            && m.getTagUUID().equals(measurement.getTagUUID())).findFirst();
            databaseMeasurement.ifPresentOrElse(dm -> {
                if (measurement.getTemperature() != null && measurement.getTemperature() > 0)
                    dm.setTemperature(measurement.getTemperature());

                if (measurement.getHumidity() != null && measurement.getHumidity() > 0)
                    dm.setHumidity(measurement.getHumidity());

                if (measurement.getVoltage() != null && measurement.getVoltage() > 0)
                    dm.setVoltage(measurement.getVoltage());

                if (measurement.getSignal() != null && measurement.getSignal() > 0)
                    dm.setSignal(measurement.getSignal());

                measurementRepository.save(dm);
            }, () -> {
                measurementRepository.save(measurement);
            });
        }
    }
}
