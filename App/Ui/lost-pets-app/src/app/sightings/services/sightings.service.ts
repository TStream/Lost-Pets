import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from "rxjs/operators";
import { IPosting, IPostingsResponse } from 'src/app/postings/models/posting';
import { ISighting, ISightingsResponse, SightingRequest } from 'src/app/sightings/models/sighting';
import { environment } from 'src/environments/environment';

@Injectable({
  providedIn: 'root'
})
export class SightingsService {
  private url: string = environment.lostPetAPI

  constructor(private readonly http: HttpClient) { }

  getAllSightings() :Observable<Array<ISighting>>{
    return this.http.get<ISightingsResponse>(`${this.url}sightings`).pipe(
      map((res: ISightingsResponse) => {
        return res.sightings
      }),
    );
  }

  getSightingByID(id: number) :Observable<ISighting>{
    return this.http.get<ISightingResponse>(`${this.url}sightings/${id}`).pipe(
      map((res: ISightingResponse) => {
        return res.sighitng
      }),
    );
  }

  getSightingByGUID(guid: string) :Observable<ISighting>{
    return this.http.get<ISightingResponse>(`${this.url}sightings/private/${guid}`).pipe(
      map((res: ISightingResponse) => {
        return res.sighitng
      }),
    );
  }

  getSightingMatches(guid: string) :Observable<Array<IPosting>> {
    return this.http.get<IPostingsResponse>(`${this.url}sightings/private/${guid}/macthes`).pipe(
      map((res: IPostingsResponse) => {
        return res.postings
      }),
    );
  }

  createSighting(sighting: SightingRequest){
    return this.http.post(`${this.url}sightings`, sighting, { observe: 'response' });
  }
}

interface ISightingResponse {
  sighitng: ISighting
}
