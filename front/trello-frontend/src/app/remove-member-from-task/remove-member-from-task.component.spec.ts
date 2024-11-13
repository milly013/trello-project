import { ComponentFixture, TestBed } from '@angular/core/testing';

import { RemoveMemberFromTaskComponent } from './remove-member-from-task.component';

describe('RemoveMemberFromTaskComponent', () => {
  let component: RemoveMemberFromTaskComponent;
  let fixture: ComponentFixture<RemoveMemberFromTaskComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [RemoveMemberFromTaskComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(RemoveMemberFromTaskComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
