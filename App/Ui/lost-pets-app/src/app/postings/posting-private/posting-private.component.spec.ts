import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PostingPrivateComponent } from './posting-private.component';

describe('PostingPrivateComponent', () => {
  let component: PostingPrivateComponent;
  let fixture: ComponentFixture<PostingPrivateComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ PostingPrivateComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(PostingPrivateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
