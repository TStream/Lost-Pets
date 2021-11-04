import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from "@angular/router";
import { ClrDatagridStateInterface } from "@clr/angular/data/datagrid";
import { Observable, of } from "rxjs";
import { catchError, map } from "rxjs/operators";
import { ISighting } from '../models/sighting';
import { SightingsService } from '../services/sightings.service';

@Component({
  selector: 'app-sightings-gird',
  templateUrl: './sightings-gird.component.html',
  styleUrls: ['./sightings-gird.component.css']
})
export class SightingsGirdComponent implements OnInit {
  loading: boolean = true;
  sightings$: Observable<Array<ISighting>>;

  constructor(private readonly sightingsService: SightingsService,
    private readonly router: Router,
    private readonly route: ActivatedRoute) { }

  ngOnInit(): void {
    this.sightings$ = this.sightingsService.getAllSightings()
  }

  onAdd(): void {
    this.router.navigate(["create"], { relativeTo: this.route });
  }

}
