import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from "rxjs/operators";
import { ISighting, ISightingsResponse } from 'src/app/sightings/models/sighting';
import { environment } from 'src/environments/environment';
import { IPosting, IPostingsResponse, PostingRequest } from '../models/posting';

@Injectable({
  providedIn: 'root'
})
export class PostingsService {
  private url: string = environment.lostPetAPI

  constructor(private readonly http: HttpClient) { }

  getAllPostings() :Observable<Array<IPosting>>{
    return this.http.get<IPostingsResponse>(`${this.url}postings`).pipe(
      map((res: IPostingsResponse) => {
        return res.postings
      }),
    );
  }

  getPostingByID(id: number) :Observable<IPosting>{
    return this.http.get<IPostingResponse>(`${this.url}postings/${id}`).pipe(
      map((res: IPostingResponse) => {
        return res.posting
      }),
    );
  }

  getPostingByGUID(guid: string) :Observable<IPosting>{
    return this.http.get<IPostingResponse>(`${this.url}postings/private/${guid}`).pipe(
      map((res: IPostingResponse) => {
        return res.posting
      }),
    );
  }

  getPostingMatches(guid: string) :Observable<Array<ISighting>> {
    return this.http.get<ISightingsResponse>(`${this.url}postings/private/${guid}/macthes`).pipe(
      map((res: ISightingsResponse) => {
        return res.sightings
      }),
    );
  }

  createPosting(posting: PostingRequest){
    return this.http.post(`${this.url}postings`, posting, { observe: 'response' });
  }
}

interface IPostingResponse {
  posting: IPosting
}
