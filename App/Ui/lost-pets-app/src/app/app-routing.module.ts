import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { HomeComponent } from './layout/home/home.component';
import { PostingFormComponent } from './postings/posting-form/posting-form.component';
import { PostingPrivateComponent } from './postings/posting-private/posting-private.component';
import { PostingComponent } from './postings/posting/posting.component';
import { PostingsGridComponent } from './postings/postings-grid/postings-grid.component';
import { PostingResolver } from './postings/resolvers/posting.resolver';
import { SightingFormComponent } from './sightings/sighting-form/sighting-form.component';
import { SightingPrivateComponent } from './sightings/sighting-private/sighting-private.component';
import { SightingComponent } from './sightings/sighting/sighting.component';
import { SightingsGirdComponent } from './sightings/sightings-gird/sightings-gird.component';

const routes: Routes = [
  {
    path: "",
    component: HomeComponent,
    children: [
      {
        path: "postings",
        children: [
          {
            path: '',
            component:PostingsGridComponent
          },
          {
            path: "create",
            component: PostingFormComponent,
          },
          {
            path: ":id",
            component: PostingComponent,
            resolve: {
              posting: PostingResolver
            }
          },
          {
            path: "private/:guid",
            component: PostingPrivateComponent
          }
        ],
      },
      {
        path: "sightings",
        children: [
          {
            path: '',
            component:SightingsGirdComponent
          },
          {
            path: ":id",
            component: SightingComponent,
          },
          {
            path: "create",
            component: SightingFormComponent,
          },
          {
            path: "private/:guid",
            component: SightingPrivateComponent
          }
        ],
      }
    ],
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
