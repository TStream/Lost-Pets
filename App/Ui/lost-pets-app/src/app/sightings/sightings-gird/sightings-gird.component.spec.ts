import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SightingsGirdComponent } from './sightings-gird.component';

describe('SightingsGirdComponent', () => {
  let component: SightingsGirdComponent;
  let fixture: ComponentFixture<SightingsGirdComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ SightingsGirdComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(SightingsGirdComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
