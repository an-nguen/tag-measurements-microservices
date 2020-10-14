import {Component, Inject, OnInit} from '@angular/core';
import {MAT_DIALOG_DATA, MatDialogRef} from "@angular/material/dialog";
import {Tag} from "../_domains/tag";

export interface TemperatureZoneSelectTagData {
  tags: Tag[];
  selectedTags?: Tag[];
}

@Component({
  selector: 'app-temperature-zone-settings-select-tags-dialog',
  templateUrl: './temperature-zone-settings-select-tags-dialog.component.html',
  styleUrls: ['./temperature-zone-settings-select-tags-dialog.component.css']
})
export class TemperatureZoneSettingsSelectTagsDialogComponent implements OnInit {
  selectedTagsModel: Tag[] = [];
  unselectedTagsModel: Tag[] = [];
  leftTags: Tag[] = [];
  rightTags: Tag[] = [];

  constructor(public dialogRef: MatDialogRef<TemperatureZoneSettingsSelectTagsDialogComponent>,
              @Inject(MAT_DIALOG_DATA) public data: TemperatureZoneSelectTagData
  ) { }

  ngOnInit(): void {
    console.log(this.data)
    this.leftTags = this.data.tags;
    this.sortTags({tags: this.leftTags});
    if (this.data.selectedTags && this.data.selectedTags.length > 0) {
      this.leftTags = this.leftTags.filter( (el) => {
        return !this.data.selectedTags.includes(el);
      });
      this.rightTags.push(...this.data.selectedTags);
    }
    this.sortTags({tags: this.rightTags});
  }

  private sortTags({tags}) {
    tags = tags.sort((a, b) => {
      if (a.name > b.name) return 1;
      if (a.name < b.name) return -1;
      return 0;
    } );
  }

  transferToRight(tag?: Tag) {
    if (!tag && this.unselectedTagsModel) {
      this.rightTags.push(...this.unselectedTagsModel);
      this.leftTags = this.leftTags.filter((el) => !this.unselectedTagsModel.includes(el));
    } else {
      this.rightTags.push(tag);
      this.leftTags = this.leftTags.filter((el) => tag !== el);
    }
    this.sortTags({tags: this.leftTags});
    this.sortTags({tags: this.rightTags});
  }

  transferToLeft(tag?: Tag) {
    if (!tag && this.selectedTagsModel) {
      this.leftTags.push(...this.selectedTagsModel);
      this.rightTags = this.rightTags.filter((el) => !this.selectedTagsModel.includes(el))
    } else {
      this.leftTags.push(tag);
      this.rightTags = this.rightTags.filter((el) => tag !== el);
    }
    this.sortTags({tags: this.leftTags});
    this.sortTags({tags: this.rightTags});
  }

  apply() {
    this.dialogRef.close({selectedTags: this.rightTags});
  }

  close() {
    this.dialogRef.close();
  }

  transfer($event) {
    console.log($event);
  }
}
