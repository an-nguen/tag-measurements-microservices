import {Component, Inject, OnInit} from '@angular/core';
import {Role} from "../_domains/role";
import {MAT_DIALOG_DATA, MatDialogRef} from "@angular/material/dialog";
import {MatSnackBar} from "@angular/material/snack-bar";
import {RoleService} from "../_services/role.service";
import {User} from "../_domains/user";

interface RoleDialogData {
  role: Role;
  mode: 'create' | 'edit' | undefined;
}

@Component({
  selector: 'app-edit-role-dialog',
  templateUrl: './edit-role-dialog.component.html',
  styleUrls: ['./edit-role-dialog.component.css']
})
export class EditRoleDialogComponent implements OnInit {
  nameValue: string = '';
  privilegeRoleList: Role[] = [];

  constructor(public dialogRef: MatDialogRef<EditRoleDialogComponent>,
              public roleService: RoleService,
              private snackBar: MatSnackBar,
              @Inject(MAT_DIALOG_DATA) public data: RoleDialogData) { }

  ngOnInit(): void {
    this.nameValue = this.data.role.name;
  }

  createRole() {
    const role = {
      id: 0,
      name: this.nameValue,
    }
    this.roleService.createRole(role)
      .subscribe((resp: Role) => {
        this.snackBar.open(`Роль ${resp.name} создана.`, 'Закрыть', {
          duration: 5000
        });
        this.dialogRef.close();
      }, error => {
        this.snackBar.open(`Роль не создана. Ошибка: ${error.error}`, 'Закрыть', {
          duration: 5000
        });
      });
  }

  editRole() {
    const role = {...this.data.role};
    role.name = this.nameValue;
    this.roleService.updateRole(role)
      .subscribe((resp: Role) => {
        this.snackBar.open(`Роль ${resp.name} обновлена.`, 'Закрыть', {
          duration: 5000
        });
        this.dialogRef.close();
      }, error => {
        this.snackBar.open(`Роль не обновлена. Ошибка: ${error.error}`, 'Закрыть', {
          duration: 5000
        });
        this.dialogRef.close();
      });
  }
}
