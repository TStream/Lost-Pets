import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SightingPrivateComponent } from './sighting-private.component';

describe('SightingPrivateComponent', () => {
  let component: SightingPrivateComponent;
  let fixture: ComponentFixture<SightingPrivateComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ SightingPrivateComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(SightingPrivateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
