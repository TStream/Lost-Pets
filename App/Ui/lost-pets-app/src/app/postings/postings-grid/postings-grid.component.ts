import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from "@angular/router";
import { ClrDatagridStateInterface } from "@clr/angular/data/datagrid";
import { Observable, of } from "rxjs";
import { catchError, map } from "rxjs/operators";
import { IPosting } from 'src/app/postings/models/posting';
import { PostingsService } from 'src/app/postings/services/postings.service';
import { environment } from 'src/environments/environment';

@Component({
  selector: 'app-postings-grid',
  templateUrl: './postings-grid.component.html',
  styleUrls: ['./postings-grid.component.css']
})
export class PostingsGridComponent implements OnInit {
  loading: boolean = true;
  postings$: Observable<Array<IPosting>>;
  pictureUrl: string = environment.lostPetAPI + "pet-picutes/"

  constructor(private readonly postingsService: PostingsService,
    private readonly router: Router,
    private readonly route: ActivatedRoute) { }

  ngOnInit(): void {
    this.postings$ = this.postingsService.getAllPostings()
    .pipe(
      map((res) => {
        this.loading = false;
        return res;
      }),
      catchError((err) => {
        this.loading = false;
        return of([]);
      })
    );
  }

  onAdd(): void {
    this.router.navigate(["create"], { relativeTo: this.route });
  }

  onView(id : number) :void {
    this.router.navigate([id], {
      relativeTo: this.route,
    });
  }

}
