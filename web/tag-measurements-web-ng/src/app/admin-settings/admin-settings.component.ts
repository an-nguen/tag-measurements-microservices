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
import {Privilege} from "../_domains/privilege";
import {EditRoleDialogComponent} from "../edit-role-dialog/edit-role-dialog.component";
import {PrivilegeService} from "../_services/privilege.service";
import {EditPrivilegeDialogComponent} from "../edit-privilege-dialog/edit-privilege-dialog.component";

@Component({
  selector: 'app-admin-settings',
  templateUrl: './admin-settings.component.html',
  styleUrls: ['./admin-settings.component.css']
})
export class AdminSettingsComponent implements OnInit {
  @ViewChild(MatSort, {static: false}) sort: MatSort;
  userDataSource: MatTableDataSource<User>;
  roleDataSource: MatTableDataSource<Role>;
  privilegeDataSource: MatTableDataSource<Privilege>
  public userDisplayedColumns = ['username', 'actions'];
  public roleDisplayedColumns = ['name', 'actions'];
  public privilegeDisplayedColumns = ['name', 'value'];
  constructor(private userService: UserService,
              private roleService: RoleService,
              private privilegeService: PrivilegeService,
              private router: Router,
              private snackBar: MatSnackBar,
              private dialog: MatDialog) { }

  ngOnInit(): void {
    if (!this.roleService.isAdmin()) {
      this.router.navigate(['']);
    }
    this.updateData();
  }


  private updateData() {
    this.userService.getUsers().add(() => {
      this.userDataSource = new MatTableDataSource<User>(this.userService.users);
      this.userDataSource.sort = this.sort;
    });
    this.roleService.getRoles().add(() => {
      this.roleDataSource = new MatTableDataSource<Role>(this.roleService.roles);
      this.roleDataSource.sort = this.sort;
    })
    this.privilegeService.getPrivileges().add(() => {
      this.privilegeDataSource = new MatTableDataSource<Privilege>(this.privilegeService.privileges);
      this.privilegeDataSource.sort = this.sort;
    })
  }

  createUser() {
    this.dialog.open(EditUserDialogComponent, {
      width: '640px',
      data: {
        user: {},
        mode: 'create'
      }
    }).afterClosed().subscribe(() => {
      this.userService.getUsers().add(() => {
        this.userDataSource.data = this.userService.users;
      })
    })
  }

  editUser(user: User) {
    this.dialog.open(EditUserDialogComponent, {
      width: '640px',
      data: {
        user: user,
        mode: 'edit'
      }
    }).afterClosed().subscribe(() => {
      this.userService.getUsers().add(() => {
        this.userDataSource.data = this.userService.users;
      })
    })
  }

  deleteUser(user: User) {
    this.userService.deleteUser(user).subscribe((resp: any) => {
      this.userService.getUsers().add(() => {
        this.userDataSource.data = this.userService.users;
      })
      this.snackBar.open(`Пользователь ${user.username} удалён.`, 'Закрыть', {
        duration: 5000
      });
    }, error => {
      this.snackBar.open(`Пользователь не удалён. Ошибка: ${error.error}`, 'Закрыть', {
        duration: 5000
      });
    });
  }

  createRole() {
    this.dialog.open(EditRoleDialogComponent, {
      width: '640px',
      data: {
        role: {},
        mode: 'create'
      }
    }).afterClosed().subscribe(() => {
      this.roleService.getRoles().add(() => {
        this.roleDataSource.data = this.roleService.roles;
      })
    })
  }

  editRole(role: Role) {
    this.dialog.open(EditRoleDialogComponent, {
      width: '640px',
      data: {
        role: role,
        mode: 'edit'
      }
    }).afterClosed().subscribe(() => {
      this.roleService.getRoles().add(() => {
        this.roleDataSource.data = this.roleService.roles;
      })
    })
  }

  deleteRole(role: Role) {
    this.roleService.deleteRole(role)
        .subscribe((resp) => {
          this.roleService.getRoles().add(() => {
            this.roleDataSource.data = this.roleService.roles;
          })
              this.snackBar.open(`Роль ${role.name} удалена.`, 'Закрыть', {
                duration: 5000
              });
            }, error => {
              this.snackBar.open(`Роль не удалена. Ошибка: ${error.error}`, 'Закрыть', {
                duration: 5000
              });
            }
        );
  }

  createPrivilege() {
    this.dialog.open(EditPrivilegeDialogComponent, {
      width: '640px',
      data: {
        privilege: {},
        mode: 'create'
      }
    }).afterClosed().subscribe(() => {
      this.privilegeService.getPrivileges().add(() => {
        this.privilegeDataSource.data = this.privilegeService.privileges;
      });
    })
  }

  editPrivilege(privilege: Privilege) {
    this.dialog.open(EditPrivilegeDialogComponent, {
      width: '640px',
      data: {
        privilege: privilege,
        mode: 'edit'
      }
    }).afterClosed().subscribe(() => {
      this.privilegeService.getPrivileges().add(() => {
        this.privilegeDataSource.data = this.privilegeService.privileges;
      });
    })
  }

  deletePrivilege(removedPrivilege: Privilege) {
    this.privilegeService.deletePrivilege(removedPrivilege)
        .subscribe((resp) => {
          this.privilegeService.getPrivileges().add(() => {
            this.privilegeDataSource.data = this.privilegeService.privileges;
          });
          this.snackBar.open(`Привилегия ${removedPrivilege.name} удалена.`, 'Закрыть', {
            duration: 5000
          });
        })
  }
}
