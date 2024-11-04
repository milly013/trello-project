import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { ProjectService, Project } from '../service/project.service';

@Component({
  selector: 'app-create-project',
  standalone: true,
  templateUrl: './create-project-component.component.html',
  styleUrls: ['./create-project-component.component.css'],
  imports: [ReactiveFormsModule], // Uverite se da je ovde sve što vam treba
  providers: [ProjectService] // Obezbeđivanje servisa
})
export class CreateProjectComponent implements OnInit {
  projectForm: FormGroup;

  constructor(private fb: FormBuilder, private projectService: ProjectService) {
    this.projectForm = this.fb.group({
      name: ['', Validators.required],
      expected_end_date: ['', Validators.required],
      min_members: [1, [Validators.required, Validators.min(1)]],
      max_members: [5, [Validators.required, Validators.min(1)]],
      manager_id: ['', Validators.required]
    });
  }

  ngOnInit(): void {}

  addProject(): void {
    if (this.projectForm.valid) {
      this.projectService.createProject(this.projectForm.value).subscribe({
        next: (response) => {
          console.log('Project added successfully', response);
          this.projectForm.reset();
        },
        error: (error) => {
          console.error('Error adding project', error);
        },
        complete: () => {
          console.log('Add project observable completed');
        }
      });
    }
  }
}
