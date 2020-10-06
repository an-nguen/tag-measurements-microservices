import {Component, OnInit, ViewChild} from '@angular/core';
import {NgForm} from "@angular/forms";
import {MatTableDataSource} from "@angular/material/table";
import {User} from "../_domains/user";
import {RoleService} from "../_services/role.service";
import {Router} from "@angular/router";
import {UserService} from "../_services/user.service";
import {MatSort} from "@angular/material/sort";
import {MatDialog} from "@angular/material/dialog";
import {EditUserDialogComponent} from "../edit-user-dialog/edit-user-dialog.component";
import {MatSnackBar} from "@angular/material/snack-bar";
import {Role} from "../_domains/role";

@Component({
  selector: 'app-admin-settings',
  templateUrl: './admin-settings.component.html',
  styleUrls: ['./admin-settings.component.css']
})
export class AdminSettingsComponent implements OnInit {
  @ViewChild(MatSort, {static: false}) sort: MatSort;
  userDataSource: MatTableDataSource<User>
  roleDataSource: MatTableDataSource<Role>
  public displayedColumns = ['username', 'actions'];
  constructor(private userService: UserService,
              private roleService: RoleService,
              private router: Router,
              private snackBar: MatSnackBar,
              private dialog: MatDialog) { }

  ngOnInit(): void {
    if (!this.roleService.isAdmin()) {
      this.router.navigate(['']);
    }

    this.userService.getUsers().add(() => {
      this.userDataSource = new MatTableDataSource<User>(this.userService.users);
      this.userDataSource.sort = this.sort;
    });
    this.roleService.getRoleByToken()
  }


  createUser() {
    this.dialog.open(EditUserDialogComponent, {
      width: '640px',
      data: {
        username: '',
        mode: 'create'
      }
    })
  }

  editUser(user: User) {
    this.dialog.open(EditUserDialogComponent, {
      width: '640px',
      data: {
        username: user.username,
        mode: 'edit'
      }
    })
  }

  deleteUser(user: User) {
    this.userService.deleteUser(user).subscribe((resp: any) => {
      this.snackBar.open(`Пользователь ${user.username} удалён.`, 'Закрыть', {
        duration: 5000
      });
    });
  }
}
