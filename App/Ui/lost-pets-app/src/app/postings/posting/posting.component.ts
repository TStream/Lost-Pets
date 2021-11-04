import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { IPosting } from '../models/posting';

@Component({
  selector: 'app-posting',
  templateUrl: './posting.component.html',
  styleUrls: ['./posting.component.css']
})
export class PostingComponent implements OnInit {
  posting: IPosting

  constructor(
    private route: ActivatedRoute
  ) { }

  ngOnInit(): void {
    this.route.data.subscribe((data) => {
      this.posting = data.posting;
    });
  }

}
