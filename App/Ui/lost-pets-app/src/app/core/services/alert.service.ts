import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Subject } from 'rxjs';
import { IAlert } from '../models/alert.interface';

@Injectable({
  providedIn: 'root'
})
export class AlertService {
  private readonly storedAlert = new Subject<IAlert>();
  alert = this.storedAlert.asObservable();

  constructor() { }


  getAlert(): Observable<IAlert> {
    return this.alert;
  }

  updatedAlert(alert: IAlert): void {
    this.storedAlert.next(alert);
  }
}
