import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SightingComponent } from './sighting.component';

describe('SightingComponent', () => {
  let component: SightingComponent;
  let fixture: ComponentFixture<SightingComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ SightingComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(SightingComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
