import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PostingFormComponent } from './posting-form.component';

describe('PostingFormComponent', () => {
  let component: PostingFormComponent;
  let fixture: ComponentFixture<PostingFormComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ PostingFormComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(PostingFormComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
