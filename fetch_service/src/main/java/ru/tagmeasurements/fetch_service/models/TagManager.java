package ru.tagmeasurements.fetch_service.models;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;
import ru.tagmeasurements.fetch_service.dtos.TagManagerResponseData;

import javax.persistence.Column;
import javax.persistence.Entity;
import javax.persistence.Id;
import javax.persistence.Transient;
import java.util.Objects;

@Data
@Entity
@NoArgsConstructor
@AllArgsConstructor
public class TagManager {
    @Id
    private String mac;
    private String name;
    private String email;
    @Transient
    private boolean notifyOffline;
    @Transient
    private String notifyOfflineEmail;
    @Transient
    private boolean allowMore;
    @Transient
    private String radioId;

    @Transient
    private String sessionId;

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        TagManager that = (TagManager) o;
        return mac.equals(that.mac) &&
                name.equals(that.name) &&
                email.equals(that.email);
    }

    @Override
    public int hashCode() {
        return Objects.hash(mac, name, email);
    }

    public void parse(TagManagerResponseData response) {
        this.mac = response.getMac();
        this.name = response.getName();
        this.notifyOffline = response.isNotifyOffline();
        this.allowMore = response.isAllowMore();
        this.notifyOfflineEmail = response.getNotifyOfflineEmail();
        this.radioId = response.getRadioId();
    }
}
