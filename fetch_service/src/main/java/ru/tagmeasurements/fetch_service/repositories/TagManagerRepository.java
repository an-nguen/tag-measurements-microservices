package ru.tagmeasurements.fetch_service.repositories;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;
import ru.tagmeasurements.fetch_service.models.TagManager;

import java.util.List;

@Repository
public interface TagManagerRepository extends JpaRepository<TagManager, String> {
    List<TagManager> findAllByEmail(String email);
}
