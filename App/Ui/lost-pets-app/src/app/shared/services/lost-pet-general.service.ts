import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from 'src/environments/environment';
import { IPetType } from '../models/pet';
@Injectable({
  providedIn: 'root'
})
export class LostPetGeneralService {
  private url: string = environment.lostPetAPI

  constructor(private readonly http: HttpClient) { }

  uploadFile(file: File, fileName: string) {
    const formData: FormData = new FormData();
    formData.append('file', file, fileName);
    return this.http.post(`${this.url}pet-picutes`, formData, { observe: 'response' });
  };

  downloadFile(fileID: string): Observable<Blob> {
    return this.http.get(`${this.url}pet-picutes/${fileID}`,{responseType: 'blob'})
  }

  getPetTypes() :Observable<Array<IPetType>> {
    return this.http.get<Array<IPetType>>(`${this.url}pet-types`)
  }
}
