import { Component, OnInit } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { ActivatedRoute, Router } from '@angular/router';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { IPetType } from 'src/app/shared/models/pet';
import { LostPetGeneralService } from 'src/app/shared/services/lost-pet-general.service';
import { CatBreedService } from 'src/app/shared/services/cat-breed.service';
import { DogBreedService } from 'src/app/shared/services/dog-breed.service';
import { PostingsService } from '../services/postings.service';
import { PostingRequest } from '../models/posting';
import { AlertService } from 'src/app/core/services/alert.service';
import { alertType } from 'src/app/core/models/alert.interface';

@Component({
  selector: 'app-posting-form',
  templateUrl: './posting-form.component.html',
  styleUrls: ['./posting-form.component.css']
})
export class PostingFormComponent implements OnInit {
  form: FormGroup;
  petTypes$: Observable<Array<IPetType>>;
  dogBreeds$: Observable<Array<string>>;
  catBreeds$: Observable<Array<string>>;
  formBreeds$: Observable<Array<string>>;

  today = new Date();

  dogTypeID: number = 0
  catTypeID: number = 0

  constructor(
    private lpService: LostPetGeneralService,
    private catService: CatBreedService,
    private dogService: DogBreedService,
    private postingService: PostingsService,
    private alertService: AlertService,
    private formBuilder: FormBuilder,
    private route: ActivatedRoute,
    private router: Router
  ) { 
    this.form = this.formBuilder.group({
      date: ['', [Validators.required]],
      location: ['', [Validators.required]],
      name: [''],
      email: ['', [Validators.required]],
      petpictureid: [1],
      petname: ['',[Validators.required]],
      petcolor: [''],
      petmarks: [''],
      pettypeid: [null,[Validators.required]],
      petbreeds: [[]],
      tagshape: [''],
      tagcolor: [''],
      tagtext: [''],
    });
  }

  ngOnInit(): void {
    this.petTypes$ = this.lpService.getPetTypes().pipe(
      map((res) => {
        res.forEach(t => {
          if (t.name == "dog") this.dogTypeID = t.id
          else if (t.name == "cat") this.catTypeID = t.id
        })
        return res
      })
    )
    this.catBreeds$ = this.catService.getBreeds()
    this.dogBreeds$ = this.dogService.getBreeds()
  }

  submit() {
    if (this.form.valid) {
      console.log(PostingRequest.adapt(this.form.value))
      this.postingService
        .createPosting( PostingRequest.adapt(this.form.value))
        .subscribe({
          error: (err) => {
              let error:string = err.message ?? "an unknown error occurred."
              this.alertService.updatedAlert({ duration: 2000, message: `Failed to create Posting: ${error}`, type: alertType.danger })
            
          },
          complete: () => {
            this.alertService.updatedAlert({ duration: 2000, message: "new Posting added", type: alertType.success })
            this.form.reset()
          }
        });
    }
  }

  cancel(){
    //go back to operations-data table
    this.router.navigate(["../"], {relativeTo: this.route});
  }

  onTypeChange() {
    this.form.controls["petbreeds"].setValue([])
    if (this.form.controls["pettypeid"].value == this.dogTypeID){
      this.formBreeds$ = this.dogBreeds$
    }
    else {
      this.formBreeds$ = this.catBreeds$
    }

  }

}
