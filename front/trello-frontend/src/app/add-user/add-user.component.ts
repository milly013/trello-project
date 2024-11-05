import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { Component } from '@angular/core';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators, FormsModule } from '@angular/forms';

@Component({
  selector: 'app-add-user',
  standalone: true, // Oznaka za standalone komponentu
  templateUrl: './add-user.component.html',
  imports: [ReactiveFormsModule,CommonModule,FormsModule] // Uvezi ReactiveFormsModule
})
export class AddUserComponent {
  userProjectForm: FormGroup;

  constructor(
    private fb: FormBuilder,
    private http: HttpClient // Dodato za HTTP zahteve
  ) {
    this.userProjectForm = this.fb.group({
      projectId: ['', Validators.required],
      userId: ['', Validators.required]
    });
  }

  addUserToProject() {
    if (this.userProjectForm.valid) {
      const projectId = this.userProjectForm.value.projectId;
      const userId = this.userProjectForm.value.userId;

      // Definiši payload za slanje
      const payload = { projectId, userId };

      // Pošalji POST zahtev na backend endpoint
      this.http.post('http://localhost:8080/users', payload)
        .subscribe({
          next: (response: any) => {
            console.log('Korisnik uspešno dodat u projekat', response);
          },
          error: (error: any) => {
            console.error('Greška prilikom dodavanja korisnika', error);
          }
        });
    }
  }
}