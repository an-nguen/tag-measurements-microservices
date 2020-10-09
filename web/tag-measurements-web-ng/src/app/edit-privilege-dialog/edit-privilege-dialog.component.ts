import {Component, Inject, OnInit} from '@angular/core';
import {MAT_DIALOG_DATA, MatDialogRef} from "@angular/material/dialog";
import {Privilege} from "../_domains/privilege";
import {PrivilegeService} from "../_services/privilege.service";
import {MatSnackBar} from "@angular/material/snack-bar";
import {Role} from "../_domains/role";

interface EditPrivilegeData {
  privilege: Privilege;
  mode: 'create' | 'edit' | undefined;
}

@Component({
  selector: 'app-edit-privilege-dialog',
  templateUrl: './edit-privilege-dialog.component.html',
  styleUrls: ['./edit-privilege-dialog.component.css']
})
export class EditPrivilegeDialogComponent implements OnInit {
  nameValue: string = '';
  valValue: string = '';
  rolesPrivilege: Role[] = [];

  constructor(public dialogRef: MatDialogRef<EditPrivilegeDialogComponent>,
              @Inject(MAT_DIALOG_DATA) public data: EditPrivilegeData,
              public privilegeService: PrivilegeService,
              private snackBar: MatSnackBar
  ) { }

  ngOnInit(): void {
    if (this.data.privilege.name) {
      this.nameValue = this.data.privilege.name;
    }
    if (this.data.privilege.value) {
      this.valValue = this.data.privilege.value;
    }
  }

  createPrivilege() {
    const privilege = {
      name: this.nameValue,
      value: this.valValue
    }
    this.privilegeService.createPrivilege(privilege)
      .subscribe((resp: Privilege) => {
        this.snackBar.open(`Привилегия ${resp.name} создана.`, 'Закрыть', {
          duration: 5000
        });
        this.dialogRef.close();
      }, error => {
        this.snackBar.open(`Привилегия не создана. Ошибка: ${error.error}`, 'Закрыть', {
          duration: 5000
        });
        this.dialogRef.close();
      });
  }

  editPrivilege() {
    const privilege = {
      id: this.data.privilege.id,
      name: this.nameValue,
      value: this.valValue,
      roles: this.data.privilege.roles
    }
    this.privilegeService.updatePrivilege(privilege).subscribe((resp: Privilege) => {
      this.snackBar.open(`Привилегия ${resp.name} обновлена.`, 'Закрыть', {
        duration: 5000
      });
      this.dialogRef.close();
    }, error => {
      this.snackBar.open(`Привилегия не обновлена. Ошибка: ${error.error}`, 'Закрыть', {
        duration: 5000
      });
      this.dialogRef.close();
    })
  }
}
