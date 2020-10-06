import {Component, Inject, OnInit} from '@angular/core';
import {MAT_DIALOG_DATA, MatDialogRef} from "@angular/material/dialog";
import {UserService} from "../_services/user.service";
import {User} from "../_domains/user";
import {MatSnackBar} from "@angular/material/snack-bar";

interface UserDialogData {
  username: string,
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

  constructor(public dialogRef: MatDialogRef<EditUserDialogComponent>,
              public userService: UserService,
              private snackBar: MatSnackBar,
              @Inject(MAT_DIALOG_DATA) public data: UserDialogData) { }

  ngOnInit(): void {
  }

  editUser() {
    if (!this.usernameValue && !this.passwordValue) {
      return
    }
    const user = {
      id: -1,
      username: this.usernameValue,
      password: this.passwordValue
    }
    this.userService.updateUser(user).subscribe((resp: User) => {
      this.userService.users.push(resp);
      this.snackBar.open(`Пользователь ${resp.username} изменён.`, 'Закрыть', {
        duration: 5000
      });
    }, error => {
      this.snackBar.open(`Пользователь не изменён. Ошибка: ${error.error}.`, 'Закрыть', {
        duration: 5000
      });
    })
  }

  createUser() {
    if (!this.usernameValue && !this.passwordValue) {
      return
    }
    const user = {
      id: -1,
      username: this.usernameValue,
      password: this.passwordValue
    }
    this.userService.createUser(user).subscribe((resp: User) => {
      this.userService.users.push(resp);
      this.snackBar.open(`Пользователь ${resp.username} добавлен.`, 'Закрыть', {
        duration: 5000
      });
    }, error => {
      this.snackBar.open(`Пользователь не создан. Ошибка: ${error.error}.`, 'Закрыть', {
        duration: 5000
      });
    })
  }
}
