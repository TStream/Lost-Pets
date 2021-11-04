import { TestBed } from '@angular/core/testing';
import {
  HttpClientTestingModule,
  HttpTestingController,
} from "@angular/common/http/testing";
import { environment } from "src/environments/environment";

import { DogBreedService } from './dog-breed.service';

//example response from api
const dogBreeds = {
    "message": {
        "affenpinscher": [],
        "basenji": [],
        "buhund": [
            "norwegian"
        ],
        "bulldog": [
            "boston",
            "english",
            "french"
        ],
    }
}


const url: string = environment.dogAPI;

describe('DogBreedService', () => {
  let httpTestingController: HttpTestingController;
  let service: DogBreedService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
    });
    service = TestBed.inject(DogBreedService);
    httpTestingController = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    // After every test, assert that there are no more pending requests.
    httpTestingController.verify();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  describe("Get Dog Breeds", () => {
    it("should return array of dog breeds as strings", () => {
      let expected: Array<string> = ["affenpinscher","basenji", "norwegian buhund","boston bulldog","english bulldog", "french bulldog"]
      service.getBreeds().subscribe(res => {
        expect(res).toEqual(expected)
      });

      const http = httpTestingController.expectOne({
        url:        `${url}api/breeds/list/all`,
        method:"GET"
      });

      http.flush(dogBreeds)
    })
  })
});
