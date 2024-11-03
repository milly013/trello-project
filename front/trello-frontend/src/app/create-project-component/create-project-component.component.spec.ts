import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CreateProjectComponentComponent } from './create-project-component.component';

describe('CreateProjectComponentComponent', () => {
  let component: CreateProjectComponentComponent;
  let fixture: ComponentFixture<CreateProjectComponentComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CreateProjectComponentComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CreateProjectComponentComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
