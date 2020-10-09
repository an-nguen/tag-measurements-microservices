package ru.tagmeasurements.fetch_service.models;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;
import org.hibernate.annotations.Type;
import org.hibernate.annotations.TypeDef;
import org.hibernate.annotations.TypeDefs;
import ru.tagmeasurements.fetch_service.dtos.TagResponseData;
import ru.tagmeasurements.fetch_service.dtos.TempSensorConfig;

import javax.persistence.*;
import java.time.LocalDate;
import java.util.List;
import java.util.Objects;
import java.util.UUID;

@Data
@Entity
@NoArgsConstructor
@AllArgsConstructor
public class Tag {
    @Id
    private UUID uuid;
    private String name;
    private String mac;
    private LocalDate verificationDate;
    private Double higherTemperatureLimit;
    private Double lowerTemperatureLimit;
    @ManyToMany
    private List<TemperatureZone> temperatureZones;

    @Transient
    private Long slaveId;
    @Transient
    private Integer tagType;

    public void parse(TagResponseData response) {
        uuid = UUID.fromString(response.getUuid());
        name = response.getName();
        slaveId = response.getSlaveId();
        tagType = response.getTagType();
    }

    public void parse(TempSensorConfig config) {
        lowerTemperatureLimit = config.getTh_low();
        higherTemperatureLimit = config.getTh_high();
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        Tag tag = (Tag) o;
        return uuid.equals(tag.uuid) &&
                name.equals(tag.name) &&
                mac.equals(tag.mac) &&
                Objects.equals(verificationDate, tag.verificationDate) &&
                Objects.equals(higherTemperatureLimit, tag.higherTemperatureLimit) &&
                Objects.equals(lowerTemperatureLimit, tag.lowerTemperatureLimit);
    }

    @Override
    public int hashCode() {
        return Objects.hash(uuid, name, mac, verificationDate, higherTemperatureLimit, lowerTemperatureLimit);
    }

}
