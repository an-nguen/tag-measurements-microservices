package ru.tagmeasurements.fetch_service.repositories;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;
import ru.tagmeasurements.fetch_service.models.Measurement;

import java.time.LocalDateTime;
import java.util.Collection;
import java.util.List;
import java.util.UUID;

@Repository
public interface MeasurementRepository extends JpaRepository<Measurement, Long> {
    Measurement getFirstByDateAndTagUUIDOrderByDateDesc(LocalDateTime date, UUID tagUUID);
    List<Measurement> getAllByDateInAndTagUUIDInOrderByDateDesc(Collection<LocalDateTime> date, Collection<UUID> tagUUID);
}
