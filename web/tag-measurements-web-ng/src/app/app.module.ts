import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';
import {HTTP_INTERCEPTORS, HttpClientModule} from '@angular/common/http';

import {AppRoutingModule} from './app-routing.module';
import {AppComponent} from './app.component';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {ServiceWorkerModule} from '@angular/service-worker';
import {environment} from '../environments/environment';
import {RouterModule} from '@angular/router';
import {TagListComponent} from './tag-list/tag-list.component';
import {LayoutModule} from '@angular/cdk/layout';
import {MatToolbarModule} from '@angular/material/toolbar';
import {MatButtonModule} from '@angular/material/button';
import {MatSidenavModule} from '@angular/material/sidenav';
import {MatIconModule} from '@angular/material/icon';
import {MatListModule} from '@angular/material/list';
import {MatCardModule} from '@angular/material/card';
import {MatFormFieldModule} from '@angular/material/form-field';
import {MatSelectModule} from '@angular/material/select';
import {TemperatureZoneSettings} from './temperature-zone-settings/temperature-zone-settings.component';
import {MatExpansionModule} from '@angular/material/expansion';
import {MatDialogModule} from '@angular/material/dialog';
import {SelectTagDialogComponent} from './select-tag-dialog/select-tag-dialog.component';
import {MatMenuModule} from '@angular/material/menu';
import {MatTabsModule} from '@angular/material/tabs';
import {MatInputModule} from '@angular/material/input';
import {PlotPageComponent} from './plot-page/plot-page.component';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {MatDatepickerModule} from '@angular/material/datepicker';
import {MAT_DATE_LOCALE, MatNativeDateModule} from '@angular/material/core';
import {SelectTagManagerDialogComponent} from './select-tag-manager-dialog/select-tag-manager-dialog.component';
import {WirelessTagAccountsPageComponent} from './wireless-tag-accounts-page/wireless-tag-accounts-page.component';
import {ErrorDialogComponent} from './error-dialog/error-dialog.component';
import {MatProgressBarModule} from '@angular/material/progress-bar';
import {MatSnackBarModule} from '@angular/material/snack-bar';
import {MatTooltipModule} from '@angular/material/tooltip';
import {MatTableModule} from '@angular/material/table';
import {MatSortModule} from '@angular/material/sort';

import * as PlotlyJS from 'plotly.js/dist/plotly.js';
import {PlotlyModule} from 'angular-plotly.js';
import locale from 'plotly.js-locales/ru';
import {LoginDialogComponent} from './login-dialog/login-dialog.component';
import {MatCheckboxModule} from '@angular/material/checkbox';
import {MatButtonToggleModule} from "@angular/material/button-toggle";
import {SelectionTypeBtnGroupComponent} from './selection-type-btn-group/selection-type-btn-group.component';
import {AuthInterceptor} from "./_interceptors/auth.interceptor";
import {ErrorInterceptor} from "./_interceptors/error.interceptor";
import {AuthGuard} from "./_guard/auth.guard";
import {EditTagDialogComponent} from './edit-tag-dialog/edit-tag-dialog.component';
import {TemperatureZoneSettingsSelectTagsDialogComponent} from './temperature-zone-settings-select-tags-dialog/temperature-zone-settings-select-tags-dialog.component';
import {MatGridListModule} from "@angular/material/grid-list";
import {WirelessTagAccountsEditDialogComponent} from './wireless-tag-accounts-edit-dialog/wireless-tag-accounts-edit-dialog.component';
import {AdminSettingsComponent} from './admin-settings/admin-settings.component';

PlotlyModule.plotlyjs = PlotlyJS;
PlotlyJS.register(locale);
PlotlyJS.setPlotConfig({locale: 'ru'});

@NgModule({
  declarations: [
    AppComponent,
    TagListComponent,
    TemperatureZoneSettings,
    SelectTagDialogComponent,
    PlotPageComponent,
    SelectTagManagerDialogComponent,
    WirelessTagAccountsPageComponent,
    ErrorDialogComponent,
    LoginDialogComponent,
    SelectionTypeBtnGroupComponent,
    EditTagDialogComponent,
    TemperatureZoneSettingsSelectTagsDialogComponent,
    WirelessTagAccountsEditDialogComponent,
    AdminSettingsComponent,
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    AppRoutingModule,
    PlotlyModule,
    BrowserAnimationsModule,
    ServiceWorkerModule.register('ngsw-worker.js', {enabled: environment.production}),
    RouterModule.forRoot([
      {path: '', component: TagListComponent, canActivate: [AuthGuard]},
      {path: 'warehouse-settings', component: TemperatureZoneSettings, canActivate: [AuthGuard]},
      {path: 'wst-accounts', component: WirelessTagAccountsPageComponent, canActivate: [AuthGuard]},
      {path: 'admin-settings', component: AdminSettingsComponent, canActivate: [AuthGuard]},
      {path: 'plot', component: PlotPageComponent, canActivate: [AuthGuard]},
      {path: 'login', component: LoginDialogComponent},
    ]),
    LayoutModule,
    MatToolbarModule,
    MatButtonModule,
    MatSidenavModule,
    MatIconModule,
    MatListModule,
    MatCardModule,
    MatFormFieldModule,
    MatSelectModule,
    MatDialogModule,
    MatExpansionModule,
    MatMenuModule,
    MatTabsModule,
    MatInputModule,
    FormsModule,
    MatDatepickerModule,
    MatNativeDateModule,
    MatProgressBarModule,
    MatSnackBarModule,
    ReactiveFormsModule,
    MatTooltipModule,
    MatTableModule,
    MatSortModule,
    MatCheckboxModule,
    MatButtonToggleModule,
    MatGridListModule,
  ],
  providers: [
    { provide: HTTP_INTERCEPTORS, useClass: AuthInterceptor, multi: true },
    { provide: HTTP_INTERCEPTORS, useClass: ErrorInterceptor, multi: true },
    {provide: MAT_DATE_LOCALE, useValue: 'ru-RU'},
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
