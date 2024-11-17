import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AddMemberToTaskComponent } from './add-member-to-task.component';

describe('AddMemberToTaskComponent', () => {
  let component: AddMemberToTaskComponent;
  let fixture: ComponentFixture<AddMemberToTaskComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AddMemberToTaskComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(AddMemberToTaskComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
