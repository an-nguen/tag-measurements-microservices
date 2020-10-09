package ru.tagmeasurements.fetch_service.models;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import javax.persistence.Entity;
import javax.persistence.Id;
import javax.persistence.ManyToMany;
import java.util.List;

@Data
@Entity
@NoArgsConstructor
@AllArgsConstructor
public class TemperatureZone {
    @Id
    private Long id;
    private String name;
    private String description;
    private Double lowerTempLimit;
    private Double higherTempLimit;
    private String notifyEmails;
    @ManyToMany
    private List<Tag> tags;
}
