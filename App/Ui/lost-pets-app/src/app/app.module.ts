//Angular Imports
import { CUSTOM_ELEMENTS_SCHEMA, NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { HttpClientModule } from '@angular/common/http';
import { ReactiveFormsModule, FormsModule } from '@angular/forms';

//Third Party
import { ClarityModule } from '@clr/angular';

//App Imports
import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { HomeComponent } from './layout/home/home.component';
import { PostingsGridComponent } from './postings/postings-grid/postings-grid.component';
import { PostingComponent } from './postings/posting/posting.component';
import { PostingFormComponent } from './postings/posting-form/posting-form.component';
import { PostingPrivateComponent } from './postings/posting-private/posting-private.component';
import { SightingComponent } from './sightings/sighting/sighting.component';
import { SightingFormComponent } from './sightings/sighting-form/sighting-form.component';
import { SightingsGirdComponent } from './sightings/sightings-gird/sightings-gird.component';
import { SightingPrivateComponent } from './sightings/sighting-private/sighting-private.component';
import { AlertComponent } from './shared/alert/alert.component';

@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    PostingsGridComponent,
    PostingComponent,
    PostingFormComponent,
    PostingPrivateComponent,
    SightingComponent,
    SightingFormComponent,
    SightingsGirdComponent,
    SightingPrivateComponent,
    AlertComponent
  ],
  imports: [
    AppRoutingModule,
    BrowserModule,
    BrowserAnimationsModule,
    ClarityModule,
    HttpClientModule,
    ReactiveFormsModule,
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
