package ru.tagmeasurements.fetch_service.models;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import javax.persistence.*;
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
    @JoinTable(
            name = "temperature_zone_tags",
            joinColumns = @JoinColumn(name = "temperature_zone_id"),
            inverseJoinColumns = @JoinColumn(name = "tag_uuid")
    )
    private List<Tag> tags;
}
