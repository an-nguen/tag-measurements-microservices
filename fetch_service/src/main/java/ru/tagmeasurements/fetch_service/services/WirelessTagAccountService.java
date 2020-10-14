package ru.tagmeasurements.fetch_service.services;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import ru.tagmeasurements.fetch_service.models.WirelessTagAccount;
import ru.tagmeasurements.fetch_service.repositories.WirelessTagAccountRepository;

import java.util.List;

@Service
public class WirelessTagAccountService {
  private final WirelessTagAccountRepository repository;

  @Autowired
  public WirelessTagAccountService(WirelessTagAccountRepository repository) {
    this.repository = repository;
  }

  public List<WirelessTagAccount> findAll() {
    return repository.findAll();
  }
}
