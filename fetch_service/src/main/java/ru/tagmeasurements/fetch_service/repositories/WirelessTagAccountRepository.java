package ru.tagmeasurements.fetch_service.repositories;


import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;
import ru.tagmeasurements.fetch_service.models.WirelessTagAccount;

@Repository
public interface WirelessTagAccountRepository extends JpaRepository<WirelessTagAccount, String> {
}
