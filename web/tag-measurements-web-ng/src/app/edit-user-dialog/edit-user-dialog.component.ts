import {Component, Inject, OnInit} from '@angular/core';
import {MAT_DIALOG_DATA, MatDialogRef} from "@angular/material/dialog";
import {UserService} from "../_services/user.service";
import {User} from "../_domains/user";
import {MatSnackBar} from "@angular/material/snack-bar";
import {RoleService} from "../_services/role.service";
import {Role} from "../_domains/role";

interface UserDialogData {
  user?: User,
  mode: 'create' | 'edit' | undefined;
}

@Component({
  selector: 'app-edit-user-dialog',
  templateUrl: './edit-user-dialog.component.html',
  styleUrls: ['./edit-user-dialog.component.css']
})
export class EditUserDialogComponent implements OnInit {
  usernameValue: string = '';
  passwordValue: string = '';
  selectedRoles: Role[];

  constructor(public dialogRef: MatDialogRef<EditUserDialogComponent>,
              public userService: UserService,
              public roleService: RoleService,
              private snackBar: MatSnackBar,
              @Inject(MAT_DIALOG_DATA) public data: UserDialogData) { }

  ngOnInit(): void {
    this.usernameValue = this.data.user.username;
    this.selectedRoles = this.roleService.roles.filter(role => {
      for (let r of this.data.user.roles) {
        if (r.id === role.id)
          return true;
      }
      return false;
    });
  }

  editUser() {
    if (!this.usernameValue && !this.passwordValue) {
      return
    }
    const user = {
      id: this.data.user.id,
      username: this.usernameValue,
      password: this.passwordValue,
      roles: this.selectedRoles
    }
    this.userService.updateUser(user).subscribe((resp: User) => {
      this.snackBar.open(`Пользователь ${resp.username} изменён.`, 'Закрыть', {
        duration: 5000
      });
      this.dialogRef.close({user: resp});
    }, error => {
      this.snackBar.open(`Пользователь не изменён. Ошибка: ${error.error}.`, 'Закрыть', {
        duration: 5000
      });
      this.dialogRef.close();
    })
  }

  createUser() {
    if (!this.usernameValue && !this.passwordValue) {
      return
    }
    const user = {
      username: this.usernameValue,
      password: this.passwordValue,
      roles: this.selectedRoles
    }
    this.userService.createUser(user).subscribe((resp: User) => {
      this.snackBar.open(`Пользователь ${resp.username} добавлен.`, 'Закрыть', {
        duration: 5000
      });
      this.dialogRef.close(resp);
    }, error => {
      this.snackBar.open(`Пользователь не создан. Ошибка: ${error.error}.`, 'Закрыть', {
        duration: 5000
      });
      this.dialogRef.close();
    })
  }
}
