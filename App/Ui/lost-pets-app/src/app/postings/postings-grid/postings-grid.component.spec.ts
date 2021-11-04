import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PostingsGridComponent } from './postings-grid.component';

describe('PostingsGridComponent', () => {
  let component: PostingsGridComponent;
  let fixture: ComponentFixture<PostingsGridComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ PostingsGridComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(PostingsGridComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
