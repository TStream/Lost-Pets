import { Component, OnInit } from '@angular/core';
import { Observable } from 'rxjs';
import { tap } from 'rxjs/operators';
import { IAlert } from 'src/app/core/models/alert.interface';
import { AlertService } from 'src/app/core/services/alert.service';

@Component({
  selector: 'app-alert',
  templateUrl: './alert.component.html',
  styleUrls: ['./alert.component.css']
})
export class AlertComponent implements OnInit {
  collapsed = true;
  closed = true;
  alert$: Observable<IAlert>;

  constructor(private readonly alertService: AlertService) {}

  ngOnInit(): void {
    this.alert$ = this.alertService.getAlert().pipe(
      tap((res) => {
        this.closed = false
        // if duration is set then close after timeout
        if(res.duration) {
          setTimeout(() => {
            this.closed = true;
          }, res.duration)
        }
      })
    )
  }
}
