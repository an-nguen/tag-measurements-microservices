package ru.tagmeasurements.fetch_service.repositories;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.stereotype.Repository;
import ru.tagmeasurements.fetch_service.models.Tag;

import java.util.List;
import java.util.UUID;

@Repository
public interface TagRepository extends JpaRepository<Tag, UUID> {
    @Query(value = "select t from Tag t where t.mac = ?1::macaddr", nativeQuery = true)
    List<Tag> findAllByMac(String mac);
}
