import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { ProjectService, Project } from '../service/project.service';
import { HttpClientModule } from '@angular/common/http';
import { AuthService } from '../service/auth.service';

@Component({
  selector: 'app-create-project',
  standalone: true,
  templateUrl: './create-project-component.component.html',
  styleUrls: ['./create-project-component.component.css'],
  imports: [ReactiveFormsModule],
  providers: [ProjectService] 
})
export class CreateProjectComponent implements OnInit {
  projectForm: FormGroup;

  constructor(private fb: FormBuilder, private projectService: ProjectService, private authService: AuthService) {
    this.projectForm = this.fb.group({
      name: ['', Validators.required],
      expected_end_date: ['', Validators.required],
      min_members: [1, [Validators.required, Validators.min(1)]],
      max_members: [5, [Validators.required, Validators.min(1)]],
    });    
  }

  ngOnInit(): void {}

  addProject(): void {
    if (this.projectForm.valid) {
      // Uzmi managerId iz localStorage-a
      const managerId = this.authService.getUserId()
      if (!managerId) {
        console.error('Manager ID is not available in local storage');
        return;
      }
  
      // Pripremanje podataka za projekat, uključujući managerId
      const projectData = {
        ...this.projectForm.value,
        endDate: new Date(this.projectForm.value.expected_end_date).toISOString(),
        minMembers: this.projectForm.value.min_members,
        maxMembers: this.projectForm.value.max_members,
        managerId: managerId
      };
      
      // Poziv servisa za kreiranje projekta
      this.projectService.createProject(projectData).subscribe({
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
