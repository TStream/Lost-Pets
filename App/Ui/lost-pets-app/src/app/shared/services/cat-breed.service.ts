import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { environment } from 'src/environments/environment';
import { Observable } from 'rxjs';
import { map } from "rxjs/operators";
@Injectable({
  providedIn: 'root'
})
export class CatBreedService {
  private url: string = environment.catAPI

  constructor(private readonly http: HttpClient) { }

  getBreeds() :Observable<Array<string>>{
    return this.http.get<Array<ICatBreed>>(`${this.url}v1/breeds`).pipe(
      map((res: Array<ICatBreed>) => {
        return res.map(c => c.name)
      }),
    );
  }
  

}

interface ICatBreed {
  name: string
}
