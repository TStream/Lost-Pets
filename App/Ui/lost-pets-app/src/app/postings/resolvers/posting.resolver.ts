import { Injectable } from '@angular/core';
import {
  Router, Resolve,
  RouterStateSnapshot,
  ActivatedRouteSnapshot
} from '@angular/router';
import { Observable, of } from 'rxjs';
import { IPosting } from '../models/posting';
import { PostingsService } from '../services/postings.service';

@Injectable({
  providedIn: 'root'
})
export class PostingResolver implements Resolve<IPosting> {
  constructor(private service: PostingsService) {}

  resolve(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<IPosting> {
    let idStr = route.paramMap.get('id')
    return this.service.getPostingByID(parseInt(idStr));
  }
}
