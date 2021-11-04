import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from "rxjs/operators";
import { environment } from 'src/environments/environment';

@Injectable({
  providedIn: 'root'
})
export class DogBreedService {
  private url: string = environment.dogAPI

  constructor(private readonly http: HttpClient) { }

  getBreeds() :Observable<Array<string>>{
    return this.http.get<IDogBreeds>(`${this.url}api/breeds/list/all`).pipe(
      map((res: IDogBreeds) => {
        let breeds = Array<string>();
        let keys = [ ...Object.keys(res.message) ];
        keys.forEach(k => {
          let subs = res.message[k]
          if (subs?.length && subs.length > 0 ){
            subs?.forEach(s => {
              breeds.push(s + " " + k)
            })
          }
          else {
            breeds.push(k)
          }
          
        })
        return breeds
      }),
    );
  }
}

interface IDogBreeds {
  message: IHash
}

interface IHash {
  [name: string] : Array<string>
}
