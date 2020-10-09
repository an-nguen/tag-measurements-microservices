package ru.tagmeasurements.fetch_service.services;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Service;
import ru.tagmeasurements.fetch_service.models.CloudHttpClient;
import ru.tagmeasurements.fetch_service.repositories.MeasurementRepository;
import ru.tagmeasurements.fetch_service.repositories.TagManagerRepository;
import ru.tagmeasurements.fetch_service.repositories.TagRepository;
import ru.tagmeasurements.fetch_service.repositories.WirelessTagAccountRepository;

@Service
public class ScheduledService {

    private final TagManagerRepository tagManagerRepository;
    private final TagRepository tagRepository;
    private final MeasurementRepository measurementRepository;
    private final WirelessTagAccountRepository wirelessTagAccountRepository;
    private final ApiWrapperService apiWrapperService;


    @Autowired
    public ScheduledService(TagManagerRepository tagManagerRepository,
                            TagRepository tagRepository,
                            MeasurementRepository measurementRepository, WirelessTagAccountRepository wirelessTagAccountRepository,
                            ApiWrapperService apiWrapperService) {
        this.tagManagerRepository = tagManagerRepository;
        this.tagRepository = tagRepository;
        this.measurementRepository = measurementRepository;
        this.wirelessTagAccountRepository = wirelessTagAccountRepository;
        this.apiWrapperService = apiWrapperService;
    }

    @Scheduled(fixedRate = 15000)
    private void clientFetch() {
        var accounts = wirelessTagAccountRepository.findAll();
        for (var account: accounts) {
            var client = new CloudHttpClient(account, apiWrapperService);
            client.init();
            client.store(tagManagerRepository, tagRepository, measurementRepository);
        }
    }

}
